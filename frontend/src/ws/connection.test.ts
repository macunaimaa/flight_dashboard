import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { WebSocketManager } from "./connection";

// --- Fake WebSocket --------------------------------------------------------

interface FakeSocket extends WebSocket {
  url: string;
  triggerOpen: () => void;
  triggerMessage: (data: string) => void;
  triggerClose: () => void;
  triggerError: () => void;
}

let instances: FakeSocket[] = [];

class FakeWebSocket {
  static instances: FakeSocket[] = [];
  readyState = 0;
  url: string;
  onopen: ((e: Event) => void) | null = null;
  onmessage: ((e: MessageEvent) => void) | null = null;
  onclose: ((e: CloseEvent) => void) | null = null;
  onerror: ((e: Event) => void) | null = null;

  constructor(url: string) {
    this.url = url;
    instances.push(this as unknown as FakeSocket);
  }
  close = vi.fn(() => {
    if (this.onclose) this.onclose(new CloseEvent("close"));
  });
  send = vi.fn();

  triggerOpen() {
    this.readyState = 1;
    this.onopen?.(new Event("open"));
  }
  triggerMessage(data: string) {
    this.onmessage?.(new MessageEvent("message", { data }));
  }
  triggerClose() {
    this.readyState = 3;
    this.onclose?.(new CloseEvent("close"));
  }
  triggerError() {
    this.onerror?.(new Event("error"));
  }
}

beforeEach(() => {
  instances = [];
  vi.stubGlobal("WebSocket", FakeWebSocket);
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
  vi.unstubAllGlobals();
});

// --- Tests -----------------------------------------------------------------

describe("WebSocketManager", () => {
  it("connect() opens a socket with the token in the query string and emits 'connecting'", () => {
    const statuses: string[] = [];
    const m = new WebSocketManager(
      "tok abc",
      () => {},
      (s) => statuses.push(s)
    );

    m.connect();

    expect(instances).toHaveLength(1);
    expect(instances[0].url).toContain("/api/v1/ws?token=tok%20abc");
    expect(statuses[0]).toBe("connecting");
  });

  it("emits 'connected' on open and dispatches parsed messages", () => {
    const received: unknown[] = [];
    const statuses: string[] = [];
    const m = new WebSocketManager(
      "t",
      (msg) => received.push(msg),
      (s) => statuses.push(s)
    );
    m.connect();
    const sock = instances[0];

    sock.triggerOpen();
    expect(statuses).toContain("connected");

    sock.triggerMessage(
      JSON.stringify({ type: "aircraft.update", payload: [] })
    );
    expect(received).toEqual([{ type: "aircraft.update", payload: [] }]);
  });

  it("ignores malformed JSON messages without throwing", () => {
    const received: unknown[] = [];
    const m = new WebSocketManager(
      "t",
      (msg) => received.push(msg),
      () => {}
    );
    m.connect();
    const sock = instances[0];

    expect(() => sock.triggerMessage("{not json")).not.toThrow();
    expect(received).toEqual([]);
  });

  it("reconnects with exponential backoff after a close, until disconnect()", () => {
    const statuses: string[] = [];
    const m = new WebSocketManager("t", () => {}, (s) => statuses.push(s));
    m.connect();
    expect(instances).toHaveLength(1);

    // First disconnect — should schedule reconnect at 1000ms.
    instances[0].triggerClose();
    vi.advanceTimersByTime(999);
    expect(instances).toHaveLength(1);
    vi.advanceTimersByTime(1);
    expect(instances).toHaveLength(2);

    // Second close — backoff doubles to 2000ms.
    instances[1].triggerClose();
    vi.advanceTimersByTime(1999);
    expect(instances).toHaveLength(2);
    vi.advanceTimersByTime(1);
    expect(instances).toHaveLength(3);

    // disconnect() stops further reconnects.
    m.disconnect();
    instances[2].triggerClose();
    vi.advanceTimersByTime(60_000);
    expect(instances).toHaveLength(3);
  });

  it("resets backoff after a successful 'open'", () => {
    const m = new WebSocketManager("t", () => {}, () => {});
    m.connect();
    instances[0].triggerClose();
    vi.advanceTimersByTime(1000);
    expect(instances).toHaveLength(2);

    // Successful open should reset backoff to 1000ms.
    instances[1].triggerOpen();
    instances[1].triggerClose();
    vi.advanceTimersByTime(999);
    expect(instances).toHaveLength(2);
    vi.advanceTimersByTime(1);
    expect(instances).toHaveLength(3);
  });

  it("onerror closes the socket", () => {
    const m = new WebSocketManager("t", () => {}, () => {});
    m.connect();
    const sock = instances[0];

    sock.triggerError();
    expect(sock.close).toHaveBeenCalled();
  });
});
