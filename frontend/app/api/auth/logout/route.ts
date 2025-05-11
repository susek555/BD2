import { authConfig } from "@/app/api/auth/[...nextauth]/route"
import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh"
import { getServerSession } from "next-auth/next"
import { NextResponse } from "next/server"

const API_URL = process.env.API_URL

export async function POST() {
  const session = await getServerSession(authConfig)
  if (!session) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 })
  }

  try {
    await fetchWithRefresh(`${API_URL}/logout`, {
      method: "POST",
      body: JSON.stringify({ refresh_token: session.user.refreshToken }),
    })
  } catch (err) {
    console.error("External logout failed:", err)
  }

  return NextResponse.json({ ok: true })
}
