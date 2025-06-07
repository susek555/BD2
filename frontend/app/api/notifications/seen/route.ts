// app/api/notifications/[id]/seen/route.ts
import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { getServerSession } from "next-auth";
import { authConfig } from "@/app/lib/authConfig";
import { NextRequest, NextResponse } from "next/server";

const API_URL = process.env.API_URL;

export async function PUT(
  req: NextRequest,
) {
  const session = await getServerSession(authConfig);
  if (!session) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
  }

  const response = await fetchWithRefresh(`${API_URL}/notification/seen`, {
    method: "PUT",
  });

  const result = await response.json();
  return NextResponse.json(result, { status: response.status });
}