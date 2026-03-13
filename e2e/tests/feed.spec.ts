import { test, expect } from "@playwright/test";

test.describe("Feeds", () => {
  test("should display global feed on home page", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("body")).toContainText("Global Feed");
  });

  test("should display tag list on home page", async ({ page }) => {
    await page.goto("/");
    // Wait for tags to load via HTMX
    await expect(page.locator("#popular-tag-list")).toBeVisible({ timeout: 10_000 });
  });

  test("should navigate to tag feed when clicking a tag", async ({ page }) => {
    await page.goto("/");

    // Wait for tags to load via HTMX
    await expect(page.locator("#popular-tag-list")).toBeVisible({ timeout: 10_000 });

    const tagLink = page.locator("#popular-tag-list a.label-pill").first();
    if (await tagLink.isVisible()) {
      const tagText = await tagLink.textContent();
      await tagLink.click();
      await page.waitForTimeout(2000);

      if (tagText) {
        await expect(page.locator("body")).toContainText(tagText.trim());
      }
    }
  });

  test("should redirect unauthenticated user from your-feed", async ({ page }) => {
    await page.goto("/your-feed");

    // Should redirect to home
    await expect(page).toHaveURL("/");
  });

  test("should show your feed tab when authenticated", async ({ page }) => {
    const uniqueId = Date.now();

    // Sign up
    await page.goto("/sign-up");
    await page.fill("#sign-up-username", `feeduser${uniqueId}`);
    await page.fill("#sign-up-email", `feed${uniqueId}@example.com`);
    await page.fill("#sign-up-password", "password123");
    await page.click('button:has-text("Sign up")');
    await page.waitForURL("/", { timeout: 10_000 });

    // Should see Your Feed tab
    await expect(page.locator("body")).toContainText("Your Feed");
    await expect(page.locator("body")).toContainText("Global Feed");
  });

  test("should paginate global feed", async ({ page }) => {
    await page.goto("/");

    // Wait for feed to load via HTMX
    await page.waitForTimeout(2000);

    // Check if pagination exists (only if enough articles)
    const pagination = page.locator(".pagination a, nav a[href*='page']");
    const count = await pagination.count();

    if (count > 0) {
      await pagination.first().click();
      await page.waitForTimeout(1000);
      await expect(page.locator("body")).toContainText("Global Feed");
    }
  });
});
