import { useEffect, useState } from 'react';
import { requireAdmin } from '../../lib/auth';

export const getServerSideProps = requireAdmin;

type User = { id: string; email: string; username: string; role: 'ADMIN'|'MEMBER'; createdAt: string; };

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string|null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    const res = await fetch('/api/admin/users');
    if (!res.ok) { setError('Failed to load users'); setLoading(false); return; }
    const data = await res.json();
    setUsers(data.users); setLoading(false);
  };

  useEffect(()=>{ fetchUsers(); }, []);

  const updateRole = async (id: string, role: 'ADMIN'|'MEMBER') => {
    const res = await fetch('/api/admin/users', {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, role })
    });
    if (res.ok) fetchUsers();
  };

  const removeUser = async (id: string) => {
    if (!confirm('Delete this user?')) return;
    const res = await fetch('/api/admin/users', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id })
    });
    if (res.ok) fetchUsers();
  };

  if (loading) return <p>Loading...</p>;
  if (error) return <p className="text-red-600">{error}</p>;

  return (
    <div className="card">
      <h1 className="text-xl font-semibold mb-4">User Management</h1>
      <table className="table">
        <thead>
          <tr>
            <th className="th">Email</th>
            <th className="th">Username</th>
            <th className="th">Role</th>
            <th className="th">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200">
          {users.map(u => (
            <tr key={u.id}>
              <td className="td">{u.email}</td>
              <td className="td">{u.username}</td>
              <td className="td">
                <select className="input" value={u.role} onChange={(e)=>updateRole(u.id, e.target.value as 'ADMIN'|'MEMBER')}>
                  <option value="MEMBER">MEMBER</option>
                  <option value="ADMIN">ADMIN</option>
                </select>
              </td>
              <td className="td">
                <button className="btn-secondary" onClick={()=>removeUser(u.id)}>Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
