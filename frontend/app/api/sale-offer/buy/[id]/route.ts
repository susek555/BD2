// app/api/sale-offer/buy/[id]/route.ts
import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { getServerSession } from "next-auth";
import { authConfig } from "@/app/lib/authConfig";
import { NextRequest, NextResponse } from "next/server";
import { API_URL } from "@/app/lib/constants";

export async function DELETE(
    req: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
  const session = await getServerSession(authConfig);
  if (!session) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
  }

  const { id } = await params;

  const response = await fetchWithRefresh(`${API_URL}/sale-offer/buy/${id}`, {
    method: "DELETE",
  });

  return response;
}

