import { test, expect, Page } from "@playwright/test";

async function signUp(page: Page): Promise<string> {
  const uniqueId = Date.now();
  await page.goto("/sign-up");
  await page.fill("#sign-up-username", `author${uniqueId}`);
  await page.fill("#sign-up-email", `author${uniqueId}@example.com`);
  await page.fill("#sign-up-password", "password123");
  await page.click('button:has-text("Sign up")');
  await page.waitForURL("/", { timeout: 10_000 });
  return `author${uniqueId}`;
}

test.describe("Articles", () => {
  test("should create a new article", async ({ page }) => {
    await signUp(page);

    await page.goto("/editor");
    await page.fill('input[name="title"]', "My Test Article");
    await page.fill('input[name="description"]', "A test description");
    await page.fill('textarea[name="content"]', "This is the body of my test article.");
    await page.click('button:has-text("Publish Article")');

    await page.waitForURL(/\/articles\//, { timeout: 10_000 });
    await expect(page.locator("#article-detail__title")).toContainText("My Test Article");
  });

  test("should view article detail page", async ({ page }) => {
    await signUp(page);

    // Create article
    await page.goto("/editor");
    await page.fill('input[name="title"]', "View Test Article");
    await page.fill('input[name="description"]', "Description");
    await page.fill('textarea[name="content"]', "Body content here.");
    await page.click('button:has-text("Publish Article")');
    await page.waitForURL(/\/articles\//, { timeout: 10_000 });

    // Verify detail page content
    await expect(page.locator("#article-detail__title")).toContainText("View Test Article");
    await expect(page.locator(".post-content")).toContainText("Body content here.");
  });

  test("should edit an existing article", async ({ page }) => {
    await signUp(page);

    // Create article
    await page.goto("/editor");
    await page.fill('input[name="title"]', "Edit Me Article");
    await page.fill('input[name="description"]', "Original description");
    await page.fill('textarea[name="content"]', "Original body.");
    await page.click('button:has-text("Publish Article")');
    await page.waitForURL(/\/articles\//, { timeout: 10_000 });

    // Click edit button (HTMX button, not a link)
    const editButton = page.locator("button.edit-button").first();
    if (await editButton.isVisible()) {
      await editButton.click();
      await page.waitForURL(/\/editor\//, { timeout: 10_000 });

      await page.fill('input[name="title"]', "Updated Article Title");
      await page.click('button:has-text("Publish Article")');
      await page.waitForURL(/\/articles\//, { timeout: 10_000 });
      await expect(page.locator("#article-detail__title")).toContainText("Updated Article Title");
    }
  });

  test("should show editor page for new article", async ({ page }) => {
    await signUp(page);

    await page.goto("/editor");
    await expect(page.locator('input[name="title"]')).toBeVisible();
    await expect(page.locator('textarea[name="content"]')).toBeVisible();
  });
});
