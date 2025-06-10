// app/api/images/upload/[id]/route.ts
import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { getServerSession } from "next-auth";
import { authConfig } from "@/app/lib/authConfig";
import { NextRequest, NextResponse } from "next/server";
import { API_URL } from "@/app/lib/constants";

export async function PATCH(
  req: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  const session = await getServerSession(authConfig);
  if (!session) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
  }

  const formData = await req.formData();
  const { id } = await params;

  if (!id) {
    return NextResponse.json({ error: "Missing ID in URL" }, { status: 400 });
  }

  const response = await fetchWithRefresh(`${API_URL}/image/${id}`, {
    method: "PATCH",
    body: formData,
  });

  const result = await response.json();
  return NextResponse.json(result, { status: response.status });
}
