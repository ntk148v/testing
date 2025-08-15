import { useEffect, useState } from 'react';
import { requireAdmin } from '../../lib/auth';
import { PrismaClient } from '@prisma/client';
import type { GetServerSideProps } from 'next';

const prisma = new PrismaClient();

export const getServerSideProps: GetServerSideProps = async (context) => {
  const authResult = await requireAdmin(context);

  if (!('props' in authResult)) {
    return authResult;
  }

  const users = await prisma.user.findMany({
    include: {
      roles: { include: { role: true } },
    },
  });

  return {
    props: {
      ...authResult.props,
      users: JSON.parse(JSON.stringify(users)),
    },
  };
};

type User = {
  id: string;
  email: string;
  username: string;
  roles: { role: { id: string; name: string } }[];
  createdAt: string;
};

export default function UsersPage({ users: initialUsers }: { users: User[] }) {
  const [users, setUsers] = useState<User[]>(initialUsers);
  const [loading, setLoading] = useState(false);

  const fetchUsers = async () => {
    try {
      const res = await fetch('/api/admin/users');
      if (res.ok) {
        const data = await res.json();
        setUsers(data.users);
      }
    } catch (error) {
      console.error('Failed to fetch users:', error);
    }
  };

  const updateRole = async (id: string, role: string) => {
    const res = await fetch('/api/admin/users', {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, role }),
    });
    if (res.ok) {
      fetchUsers();
    }
  };

  const removeUser = async (id: string) => {
    if (!confirm('Delete this user?')) return;
    const res = await fetch('/api/admin/users', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id }),
    });
    if (res.ok) {
      setUsers(users.filter(u => u.id !== id));
    }
  };

  return (
    <div className="container py-4">
      <h1 className="text-2xl font-bold mb-4">Admin - Users</h1>
      <div className="card">
        <table className="w-full text-left">
          <thead>
            <tr>
              <th className="th">Email</th>
              <th className="th">Username</th>
              <th className="th">Roles</th>
              <th className="th">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {users.map(u => (
              <tr key={u.id}>
                <td className="td">{u.email}</td>
                <td className="td">{u.username}</td>
                <td className="td">{u.roles.map(r => r.role.name).join(', ')}</td>
                <td className="td">
                  <select
                    onChange={(e) => updateRole(u.id, e.target.value)}
                    className="p-1 border rounded mr-2"
                    defaultValue=""
                  >
                    <option value="" disabled>Change role</option>
                    <option value="admin">Admin</option>
                    <option value="member">Member</option>
                  </select>
                  <button onClick={() => removeUser(u.id)} className="btn-danger">
                    Remove
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
