import Link from 'next/link';
import { useSession, signOut } from 'next-auth/react';
export default function NavBar() {
  const { data: session } = useSession();
  const roles = (session as any)?.user?.roles;
  return (
    <nav className="bg-white border-b">
      <div className="container flex items-center justify-between py-3">
        <div className="flex gap-4 items-center">
          <Link href="/" className="font-semibold">RBAC Starter</Link>
          {roles?.includes('admin') && <Link href="/admin/users" className="text-sm text-gray-700 hover:underline">Users</Link>}
        </div>
        <div className="flex items-center gap-3">
          {!session ? (
            <div className="flex gap-2">
              <Link href="/login" className="btn-secondary">Login</Link>
              <Link href="/signup" className="btn-primary">Sign up</Link>
            </div>
          ) : (
            <div className="flex gap-2 items-center">
              <span className="text-sm text-gray-600">Hi, {(session.user as any)?.email}</span>
              <button className="btn-secondary" onClick={()=>signOut({ callbackUrl: '/' })}>Logout</button>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
}
