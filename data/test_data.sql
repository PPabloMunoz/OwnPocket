INSERT INTO users (id, username, password_hash, email, created_at, updated_at) VALUES
(1, 'test', '$2a$10$3grmwiVsi6TxpuStxr1VieHYH0gnuEKd5ss52oH1qtwWLpeYCAEp.', null, '2026-01-01', '2026-01-01' );

-- Currencies (seeded automatically, ID 1 = EUR)
INSERT INTO currencies (code, name, symbol, decimal_places) VALUES
('EUR', 'Euro', '€', 2),
('USD', 'US Dollar', '$', 2);

-- Accounts
INSERT INTO accounts (user_id, name, type, balance, currency_id, description, is_active, created_at, updated_at) VALUES
(1, 'Main Checking', 'checking', 125000, 1, 'Primary checking account', 1, '2026-01-01', '2026-01-01'),
(1, 'Savings', 'savings', 500000, 1, 'Emergency fund', 1, '2026-01-01', '2026-01-01'),
(1, 'Credit Card', 'credit_card', -45000, 1, 'Visa Platinum', 1, '2026-01-01', '2026-01-01'),
(1, 'Cash Wallet', 'cash', 15000, 1, 'Physical wallet', 1, '2026-01-01', '2026-01-01');

-- Categories
INSERT INTO categories (user_id, name, type, color, icon, created_at) VALUES
(1, 'Salary', 'income', '#22c55e', 'briefcase', '2026-01-01'),
(1, 'Freelance', 'income', '#3b82f6', 'laptop', '2026-01-01'),
(1, 'Groceries', 'expense', '#ef4444', 'shopping-cart', '2026-01-01'),
(1, 'Rent', 'expense', '#f97316', 'home', '2026-01-01'),
(1, 'Transport', 'expense', '#a855f7', 'car', '2026-01-01'),
(1, 'Dining Out', 'expense', '#ec4899', 'utensils', '2026-01-01'),
(1, 'Subscriptions', 'expense', '#6366f1', 'tv', '2026-01-01');

-- Transactions (amount in cents)
INSERT INTO transactions (user_id, account_id, category_id, amount, type, date, description, created_at, updated_at) VALUES
(1, 1, 1, 300000, 'income', '2026-05-01', 'Monthly salary', '2026-05-01', '2026-05-01'),
(1, 1, NULL, 200000, 'transfer', '2026-05-02', 'Transfer to savings', '2026-05-02', '2026-05-02'),
(1, 1, 3, 8500, 'expense', '2026-05-03', 'Weekly groceries', '2026-05-03', '2026-05-03'),
(1, 3, 3, 3200, 'expense', '2026-05-04', 'Supermarket', '2026-05-04', '2026-05-04'),
(1, 1, 4, 120000, 'expense', '2026-05-01', 'Monthly rent', '2026-05-01', '2026-05-01'),
(1, 1, 5, 3500, 'expense', '2026-05-05', 'Gas', '2026-05-05', '2026-05-05'),
(1, 3, 6, 4500, 'expense', '2026-05-06', 'Dinner', '2026-05-06', '2026-05-06'),
(1, 1, 7, 1599, 'expense', '2026-05-07', 'Netflix', '2026-05-07', '2026-05-07'),
(1, 1, 7, 1199, 'expense', '2026-05-07', 'Spotify', '2026-05-07', '2026-05-07'),
(1, 4, 5, 2500, 'expense', '2026-05-08', 'Bus pass top-up', '2026-05-08', '2026-05-08'),
(1, 1, 3, 4200, 'expense', '2026-05-10', 'Farmer market', '2026-05-10', '2026-05-10');

-- Budgets
INSERT INTO budgets (user_id, category_id, period, amount, created_at) VALUES
(1, 3, '2026-05', 40000, '2026-05-01'),
(1, 5, '2026-05', 15000, '2026-05-01'),
(1, 6, '2026-05', 20000, '2026-05-01'),
(1, 7, '2026-05', 10000, '2026-05-01');
