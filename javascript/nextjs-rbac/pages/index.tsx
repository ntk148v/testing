import { useSession } from 'next-auth/react';
export default function Home() {
  const { data: session } = useSession();
  const role = (session as any)?.user?.role;
  return (
    <div className="card">
      <h1 className="text-2xl font-semibold mb-2">Welcome 👋</h1>
      <p className="text-gray-700">This is a minimal Next.js + NextAuth + Prisma starter with RBAC.</p>
      <ul className="list-disc ml-6 mt-3 text-gray-700">
        <li>Sign up as Member</li>
        <li>Login with email or username</li>
        <li>Admin can manage users at <code>/admin/users</code></li>
      </ul>
      {role && <p className="mt-3 text-sm text-gray-600">Your role: {role}</p>}
    </div>
  );
}
