import type { NextApiRequest, NextApiResponse } from 'next';
import { getServerSession } from 'next-auth';
import { authOptions } from '../auth/[...nextauth]';
import { PrismaClient } from '@prisma/client';
const prisma = new PrismaClient();
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const session = await getServerSession(req, res, authOptions);
  if (!session || (session.user as any)?.role !== 'ADMIN') return res.status(403).json({ message: 'Forbidden' });
  if (req.method === 'GET') {
    const users = await prisma.user.findMany({ select: { id: true, email: true, username: true, role: true, createdAt: true } });
    return res.json({ users });
  }
  if (req.method === 'PATCH') {
    const { id, role } = req.body || {};
    if (!id || !['ADMIN','MEMBER'].includes(role)) return res.status(400).json({ message: 'Invalid payload' });
    await prisma.user.update({ where: { id }, data: { role } });
    return res.json({ ok: true });
  }
  if (req.method === 'DELETE') {
    const { id } = req.body || {};
    if (!id) return res.status(400).json({ message: 'Invalid payload' });
    await prisma.user.delete({ where: { id } });
    return res.json({ ok: true });
  }
  return res.status(405).json({ message: 'Method Not Allowed' });
}
