import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  timeout: 30_000,
  retries: 1,
  use: {
    baseURL: "http://localhost:8181",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },
  webServer: {
    command: "cd .. && go run main.go",
    url: "http://localhost:8181",
    timeout: 30_000,
    reuseExistingServer: !process.env.CI,
  },
});
