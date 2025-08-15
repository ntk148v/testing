import type { AppProps } from 'next/app';
import { SessionProvider } from 'next-auth/react';
import NavBar from '../components/NavBar';
import '../styles/globals.css';
export default function MyApp({ Component, pageProps: { session, ...pageProps } }: AppProps) {
  return (
    <SessionProvider session={session}>
      <NavBar />
      <main className="container py-6">
        <Component {...pageProps} />
      </main>
    </SessionProvider>
  );
}
