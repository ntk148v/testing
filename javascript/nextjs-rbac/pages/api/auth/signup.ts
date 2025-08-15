import { PrismaClient } from "@prisma/client";
import bcrypt from "bcrypt";
import type { NextApiRequest, NextApiResponse } from "next";
const prisma = new PrismaClient();
export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== "POST")
    return res.status(405).json({ message: "Method Not Allowed" });
  const { username, email, password } = req.body || {};
  if (!username || !email || !password)
    return res.status(400).json({ message: "Missing fields" });
  try {
    const memberRole = await prisma.role.findUnique({ where: { name: "member" } });
    if (!memberRole) {
      return res.status(500).json({ message: "'member' role not found" });
    }

    const hashed = await bcrypt.hash(password, 10);
    const user = await prisma.user.create({
      data: {
        username,
        email,
        password: hashed,
        roles: {
          create: [
            {
              role: {
                connect: { id: memberRole.id },
              },
            },
          ],
        },
      },
      include: {
        roles: {
          include: {
            role: true,
          },
        },
      },
    });
    return res
      .status(201)
      .json({
        message: "User created",
        user: {
          id: user.id,
          email: user.email,
          username: user.username,
          roles: user.roles.map((userRole) => userRole.role.name),
        },
      });
  } catch (e: any) {
    return res
      .status(400)
      .json({ message: "User already exists or invalid data" });
  }
}
