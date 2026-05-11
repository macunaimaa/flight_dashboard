import { WS_URL } from "../utils/constants";
import type { WSIncomingMessage } from "./types";

type MessageHandler = (msg: WSIncomingMessage) => void;
type StatusHandler = (status: "connecting" | "connected" | "disconnected") => void;

export class WebSocketManager {
  private ws: WebSocket | null = null;
  private token: string;
  private onMessage: MessageHandler;
  private onStatus: StatusHandler;
  private reconnectDelay = 1000;
  private maxReconnectDelay = 30000;
  private shouldReconnect = true;

  constructor(token: string, onMessage: MessageHandler, onStatus: StatusHandler) {
    this.token = token;
    this.onMessage = onMessage;
    this.onStatus = onStatus;
  }

  connect() {
    this.shouldReconnect = true;
    this.onStatus("connecting");

    const url = `${WS_URL}/api/v1/ws?token=${encodeURIComponent(this.token)}`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      this.onStatus("connected");
      this.reconnectDelay = 1000;
    };

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as WSIncomingMessage;
        this.onMessage(msg);
      } catch {
        // Ignore malformed messages
      }
    };

    this.ws.onclose = () => {
      this.onStatus("disconnected");
      if (this.shouldReconnect) {
        setTimeout(() => this.connect(), this.reconnectDelay);
        this.reconnectDelay = Math.min(
          this.reconnectDelay * 2,
          this.maxReconnectDelay
        );
      }
    };

    this.ws.onerror = () => {
      this.ws?.close();
    };
  }

  disconnect() {
    this.shouldReconnect = false;
    this.ws?.close();
    this.ws = null;
  }
}
