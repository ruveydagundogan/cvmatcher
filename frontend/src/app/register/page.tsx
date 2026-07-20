"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function RegisterPage() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const register = async () => {
    const res = await fetch("http://localhost:8080/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name,
        email,
        password,
      }),
    });

    const data = await res.json();

    if (data.success) {
      alert("Registration successful");
      router.push("/");
    } else {
      alert(data.message);
    }
  };

  return (
    <main className="min-h-screen flex justify-center items-center">
      <div className="w-96 space-y-4 border p-8 rounded-xl">

        <h1 className="text-3xl font-bold">
          Create Account
        </h1>

        <input
          placeholder="Name"
          className="border p-3 rounded w-full"
          value={name}
          onChange={(e)=>setName(e.target.value)}
        />

        <input
          placeholder="Email"
          className="border p-3 rounded w-full"
          value={email}
          onChange={(e)=>setEmail(e.target.value)}
        />

        <input
          type="password"
          placeholder="Password"
          className="border p-3 rounded w-full"
          value={password}
          onChange={(e)=>setPassword(e.target.value)}
        />

        <button
          onClick={register}
          className="w-full bg-blue-600 text-white p-3 rounded"
        >
          Register
        </button>

      </div>
    </main>
  );
}