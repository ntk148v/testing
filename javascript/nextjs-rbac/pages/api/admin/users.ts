import type { NextApiRequest, NextApiResponse } from "next";
import { getServerSession } from "next-auth";
import { authOptions } from "../auth/[...nextauth]";
import { PrismaClient } from "@prisma/client";
const prisma = new PrismaClient();

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const session = await getServerSession(req, res, authOptions);
  if (!session || !(session.user as any)?.roles?.includes("admin")) {
    return res.status(403).json({ message: "Forbidden" });
  }
  if (req.method === "GET") {
    const users = await prisma.user.findMany({
      select: {
        id: true,
        email: true,
        username: true,
        roles: true,
        createdAt: true,
      },
    });
    return res.json({ users });
  }
  if (req.method === "PATCH") {
    const { id, role } = req.body || {};
    if (!id || !["admin", "member"].includes(role)) {
      return res.status(400).json({ message: "Invalid payload" });
    }

    const roleRecord = await prisma.role.findUnique({ where: { name: role } });
    if (!roleRecord) {
      return res.status(404).json({ message: "Role not found" });
    }

    // This is a simple implementation that replaces all roles with the new one.
    // A more complex app might want to add/remove roles individually.
    await prisma.userRole.deleteMany({ where: { userId: id } });
    await prisma.userRole.create({ data: { userId: id, roleId: roleRecord.id } });
    return res.json({ ok: true });
  }
  if (req.method === "DELETE") {
    const { id } = req.body || {};
    if (!id) return res.status(400).json({ message: "Invalid payload" });
    await prisma.user.delete({ where: { id } });
    return res.json({ ok: true });
  }
  return res.status(405).json({ message: "Method Not Allowed" });
}
