import { useState, useCallback } from "react";
import { searchNews } from "../../api/news";
import type { NewsArticle } from "../../api/news";

interface Props {
  onClose: () => void;
}

export function NewsPanel({ onClose }: Props) {
  const [query, setQuery] = useState("aviation");
  const [articles, setArticles] = useState<NewsArticle[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searched, setSearched] = useState(false);

  const handleSearch = useCallback(async () => {
    if (!query.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const result = await searchNews(query);
      setArticles(result.articles);
      setSearched(true);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Search failed");
    } finally {
      setLoading(false);
    }
  }, [query]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") handleSearch();
  };

  return (
    <div style={styles.panel}>
      <div style={styles.header}>
        <div style={styles.headerLeft}>
          <span style={styles.title}>Intelligence Feed</span>
        </div>
        <button onClick={onClose} style={styles.closeBtn}>&times;</button>
      </div>

      <div style={styles.searchBar}>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Search news..."
          style={styles.input}
        />
        <button onClick={handleSearch} style={styles.searchBtn} disabled={loading}>
          {loading ? "..." : "\u2315"}
        </button>
      </div>

      <div style={styles.quickTags}>
        {["aviation", "airline", "airport", "air traffic", "flight delay"].map((tag) => (
          <button
            key={tag}
            onClick={() => { setQuery(tag); }}
            style={{
              ...styles.tag,
              ...(query === tag ? styles.tagActive : {}),
            }}
          >
            {tag}
          </button>
        ))}
      </div>

      {error && <div style={styles.error}>{error}</div>}

      {!searched && !loading && (
        <div style={styles.placeholder}>Enter a search query to find news articles</div>
      )}

      <div style={styles.list}>
        {articles.map((article, i) => (
          <a
            key={i}
            href={article.url}
            target="_blank"
            rel="noopener noreferrer"
            style={styles.card}
          >
            {article.imageUrl && (
              <div style={styles.imageWrap}>
                <img src={article.imageUrl} alt="" style={styles.image} />
              </div>
            )}
            <div style={styles.cardContent}>
              <div style={styles.cardTitle}>{article.title}</div>
              <div style={styles.cardDesc}>{article.description}</div>
              <div style={styles.cardMeta}>
                <span style={styles.cardSource}>{article.source}</span>
                <span style={styles.cardDate}>{formatDate(article.publishedAt)}</span>
              </div>
            </div>
          </a>
        ))}

        {searched && articles.length === 0 && !loading && (
          <div style={styles.noResults}>No articles found</div>
        )}
      </div>
    </div>
  );
}

function formatDate(iso: string): string {
  try {
    const d = new Date(iso);
    const now = new Date();
    const diff = now.getTime() - d.getTime();
    const hours = Math.floor(diff / 3600000);
    if (hours < 1) return "just now";
    if (hours < 24) return `${hours}h ago`;
    return d.toLocaleDateString();
  } catch {
    return "";
  }
}

const styles: Record<string, React.CSSProperties> = {
  panel: {
    position: "absolute",
    top: 56,
    right: 0,
    width: 340,
    height: "calc(100% - 56px)",
    background: "#111827",
    borderLeft: "1px solid #1f2937",
    zIndex: 70,
    display: "flex",
    flexDirection: "column",
    overflow: "hidden",
  },
  header: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    padding: "12px 16px",
    borderBottom: "1px solid #1f2937",
  },
  headerLeft: {
    display: "flex",
    alignItems: "center",
    gap: 8,
  },
  title: {
    color: "#e0e6ed",
    fontWeight: 700,
    fontSize: 13,
    textTransform: "uppercase" as const,
    letterSpacing: 1,
  },
  closeBtn: {
    background: "none",
    border: "none",
    color: "#6b7280",
    fontSize: 20,
    cursor: "pointer",
    padding: "0 4px",
    lineHeight: 1,
  },
  searchBar: {
    display: "flex",
    padding: "8px 16px",
    gap: 8,
    borderBottom: "1px solid #1f2937",
  },
  input: {
    flex: 1,
    background: "#1f2937",
    border: "1px solid #374151",
    borderRadius: 4,
    padding: "8px 12px",
    color: "#e0e6ed",
    fontSize: 13,
    outline: "none",
    fontFamily: "monospace",
  },
  searchBtn: {
    background: "#2563eb",
    border: "none",
    borderRadius: 4,
    color: "white",
    padding: "8px 14px",
    fontSize: 16,
    cursor: "pointer",
  },
  quickTags: {
    display: "flex",
    flexWrap: "wrap" as const,
    gap: 6,
    padding: "8px 16px",
    borderBottom: "1px solid #1f2937",
  },
  tag: {
    background: "#1f2937",
    border: "1px solid #374151",
    borderRadius: 3,
    color: "#9ca3af",
    padding: "3px 8px",
    fontSize: 11,
    cursor: "pointer",
  },
  tagActive: {
    background: "#2563eb33",
    borderColor: "#2563eb",
    color: "#93c5fd",
  },
  list: {
    flex: 1,
    overflow: "auto",
    padding: "8px 0",
  },
  card: {
    display: "block",
    padding: "10px 16px",
    borderBottom: "1px solid #1f293744",
    textDecoration: "none",
    cursor: "pointer",
    transition: "background 0.15s",
  },
  imageWrap: {
    marginBottom: 8,
    borderRadius: 4,
    overflow: "hidden",
    maxHeight: 120,
  },
  image: {
    width: "100%",
    height: "auto",
    maxHeight: 120,
    objectFit: "cover" as const,
    display: "block",
  },
  cardContent: {},
  cardTitle: {
    color: "#e0e6ed",
    fontSize: 13,
    fontWeight: 600,
    lineHeight: 1.4,
    marginBottom: 4,
  },
  cardDesc: {
    color: "#9ca3af",
    fontSize: 11,
    lineHeight: 1.4,
    marginBottom: 6,
    display: "-webkit-box",
    WebkitLineClamp: 2,
    WebkitBoxOrient: "vertical" as const,
    overflow: "hidden",
  },
  cardMeta: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },
  cardSource: {
    color: "#2563eb",
    fontSize: 10,
    fontWeight: 600,
    textTransform: "uppercase" as const,
  },
  cardDate: {
    color: "#6b7280",
    fontSize: 10,
  },
  error: {
    color: "#ef4444",
    fontSize: 12,
    padding: "12px 16px",
    textAlign: "center" as const,
  },
  placeholder: {
    color: "#6b7280",
    fontSize: 13,
    textAlign: "center" as const,
    padding: 32,
  },
  noResults: {
    color: "#6b7280",
    fontSize: 13,
    textAlign: "center" as const,
    padding: 24,
  },
};
