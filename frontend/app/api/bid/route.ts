// app/api/bid/route.ts
import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { getServerSession } from "next-auth";
import { authConfig } from "@/app/lib/authConfig";
import { NextRequest, NextResponse } from "next/server";
import { API_URL } from "@/app/lib/constants";

export async function POST(
    req: NextRequest,
) {
  const session = await getServerSession(authConfig);
  if (!session) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
  }

  const data = await req.json();

  const response = await fetchWithRefresh(`${API_URL}/bid`, {
    method: "POST",
    body: JSON.stringify(data),
  });

  return response;
}

