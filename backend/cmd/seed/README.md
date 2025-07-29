# Kyooar Database Seeder

This seeder creates a default test user and sample data for development and testing purposes.

## Usage

### Create default user (first time)

```bash
make seed
```

### Force recreate user (if already exists)

```bash
make seed-force
```

## What gets created:

### 👤 Test User Account

- **Email:** `admin@kyooar.com`
- **Password:** `admin123`
- **Company:** `Kyooar Demo Organization`
- **Status:** Active & Email Verified
- **Subscription:** Starter Plan (if available)

### 🏢 Sample Organization

- **Name:** Demo Organization
- **Description:** A sample organization for testing Kyooar
- **Email:** admin@kyooar.com
- **Phone:** +1-555-0123
- **Website:** https://demo.kyooar.com

### 📍 Sample Location

- **Name:** Main Location
- **Address:** 123 Organization St, Food City, CA 12345, USA

### 🔳 Sample QR Code

- **Code:** `DEMO001`
- **Label:** Table 1
- **Type:** table
- **Valid for:** 1 year

### 🍽️ Sample Products

1. **Classic Burger** - $12.99 (Mains)
2. **Caesar Salad** - $8.99 (Salads)
3. **Chocolate Cake** - $6.99 (Desserts)

## Notes

- The seeder is completely optional and separate from migrations
- Use `--force` flag to recreate existing data
- All data is created with test-friendly values
- Perfect for development and demo purposes

## Frontend Testing

After running the seeder, you can immediately test the frontend authentication:

1. Visit: `http://localhost:5173/login`
2. Login with: `admin@kyooar.com` / `admin123`
3. Explore the dashboard with sample data!
