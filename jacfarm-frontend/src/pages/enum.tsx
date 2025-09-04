export const Page = {
  FLAGS_PAGE: "flags",
  EXPLOITS_PAGE: "exploits",
  TEAMS_PAGE: "teams",
  CONFIG_PAGE: "config",
  LOGS_PAGE: "logs",
} as const;
export type PageType = typeof Page[keyof typeof Page];