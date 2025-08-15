import { PrismaClient } from '@prisma/client';
import bcrypt from 'bcrypt';
const prisma = new PrismaClient();
async function main() {
  const password = await bcrypt.hash('admin123', 10);
  await prisma.user.upsert({
    where: { email: 'admin@example.com' },
    update: {},
    create: { name: 'Admin', email: 'admin@example.com', username: 'admin', password, role: 'ADMIN' }
  });
}
main().then(()=>console.log('Seed complete')).catch(e=>{console.error(e); process.exit(1)}).finally(async()=>{await prisma.$disconnect()});
