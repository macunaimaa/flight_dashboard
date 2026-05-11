import { apiClient } from "./client";

export interface NewsArticle {
  title: string;
  description: string;
  url: string;
  source: string;
  imageUrl: string;
  publishedAt: string;
}

export interface NewsSearchResult {
  articles: NewsArticle[];
  total: number;
}

export async function searchNews(query: string, lang = "en"): Promise<NewsSearchResult> {
  const params = new URLSearchParams({ q: query, lang });
  const resp = await apiClient.get<NewsSearchResult>(`/api/v1/news/search?${params}`);
  return resp.data;
}
