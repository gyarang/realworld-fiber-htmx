import { test, expect } from "@playwright/test";

test.describe("Authentication", () => {
  test("should display sign up page", async ({ page }) => {
    await page.goto("/sign-up");
    await expect(page.locator("h1")).toContainText("Sign up");
  });

  test("should display sign in page", async ({ page }) => {
    await page.goto("/sign-in");
    await expect(page.locator("h1")).toContainText("Sign in");
  });

  test("should show error on empty sign in", async ({ page }) => {
    await page.goto("/sign-in");
    await page.click('button:has-text("Sign in")');

    // HTMX swaps error messages into the form container
    await expect(
      page.locator("#sign-in-form-messages .alert-danger")
    ).toBeVisible({ timeout: 10_000 });
  });

  test("should show error on invalid credentials", async ({ page }) => {
    await page.goto("/sign-in");
    await page.fill("#sign-in-email", "nonexistent@example.com");
    await page.fill("#sign-in-password", "wrongpassword");
    await page.click('button:has-text("Sign in")');

    await expect(
      page.locator("#sign-in-form-messages .alert-danger")
    ).toBeVisible({ timeout: 10_000 });
  });

  test("should sign up a new user", async ({ page }) => {
    const uniqueId = Date.now();
    await page.goto("/sign-up");

    await page.fill("#sign-up-username", `testuser${uniqueId}`);
    await page.fill("#sign-up-email", `test${uniqueId}@example.com`);
    await page.fill("#sign-up-password", "password123");
    await page.click('button:has-text("Sign up")');

    // HTMX redirect after successful sign-up
    await page.waitForURL("/", { timeout: 10_000 });
    await expect(page.locator("nav.navbar")).toContainText("Settings");
  });

  test("should sign in and sign out", async ({ page }) => {
    const uniqueId = Date.now();

    // Sign up first
    await page.goto("/sign-up");
    await page.fill("#sign-up-username", `loginuser${uniqueId}`);
    await page.fill("#sign-up-email", `login${uniqueId}@example.com`);
    await page.fill("#sign-up-password", "password123");
    await page.click('button:has-text("Sign up")');
    await page.waitForURL("/", { timeout: 10_000 });

    // Go to settings to sign out
    await page.goto("/settings");
    await page.click('button:has-text("Or click here to logout")');
    await page.waitForURL("/", { timeout: 10_000 });

    // Sign in
    await page.goto("/sign-in");
    await page.fill("#sign-in-email", `login${uniqueId}@example.com`);
    await page.fill("#sign-in-password", "password123");
    await page.click('button:has-text("Sign in")');

    await page.waitForURL("/", { timeout: 10_000 });
    await expect(page.locator("nav.navbar")).toContainText("Settings");
  });
});
