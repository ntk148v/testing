# Next.js RBAC Starter
Features:
- Signup/Login with username or email + password (NextAuth Credentials)
- RBAC (Admin, Member) in JWT/session
- Admin-only User Management (list, change role, delete)
- Prisma + PostgreSQL
- Tailwind CSS minimal UI
## Quickstart
```bash
cp .env.example .env
# Edit DATABASE_URL and NEXTAUTH_SECRET
npm install
npx prisma migrate dev --name init
npm run seed
npm run dev
```
Admin credentials (from seed):
- Email: `admin@example.com`
- Username: `admin`
- Password: `admin123`
Visit:
- `/signup` to register a member
- `/login` to sign in
- `/admin/users` for user management (admin only)
Notes:
- Admin page protected via `getServerSideProps` (see `lib/auth.ts`).
- Admin APIs check the session role.
