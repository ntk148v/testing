import { PrismaClient } from '@prisma/client';
import bcrypt from 'bcrypt';
import type { NextApiRequest, NextApiResponse } from 'next';
const prisma = new PrismaClient();
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method !== 'POST') return res.status(405).json({ message: 'Method Not Allowed' });
  const { username, email, password } = req.body || {};
  if (!username || !email || !password) return res.status(400).json({ message: 'Missing fields' });
  try {
    const hashed = await bcrypt.hash(password, 10);
    const user = await prisma.user.create({ data: { username, email, password: hashed, role: 'MEMBER' } });
    return res.status(201).json({ message: 'User created', user: { id: user.id, email: user.email, username: user.username, role: user.role } });
  } catch (e: any) {
    return res.status(400).json({ message: 'User already exists or invalid data' });
  }
}
