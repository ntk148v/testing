import { FormEvent, useState } from "react";
import { signIn } from "next-auth/react";
import { useRouter } from "next/router";

export default function LoginPage() {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    const res = await signIn("credentials", {
      redirect: false,
      identifier,
      password,
    });
    if (res?.error) {
      setError("Invalid credentials");
      return;
    }
    router.push("/");
  };
  return (
    <div className="card max-w-md mx-auto">
      <h1 className="text-xl font-semibold mb-4">Login</h1>
      <form className="space-y-4" onSubmit={onSubmit}>
        <div>
          <label className="label">Username or Email</label>
          <input
            className="input"
            value={identifier}
            onChange={(e) => setIdentifier(e.target.value)}
          />
        </div>
        <div>
          <label className="label">Password</label>
          <input
            type="password"
            className="input"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        {error && <p className="text-red-600 text-sm">{error}</p>}
        <button className="btn-primary w-full" type="submit">
          Login
        </button>
      </form>
    </div>
  );
}
