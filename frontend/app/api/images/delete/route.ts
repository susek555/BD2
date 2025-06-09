// app/api/images/upload/[id]/route.ts
import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { getServerSession } from "next-auth";
import { authConfig } from "@/app/lib/authConfig";
import { NextRequest, NextResponse } from "next/server";
import { API_URL } from "@/app/lib/constants";

export async function DELETE(
  req: NextRequest
) {
  const session = await getServerSession(authConfig);
  if (!session) {
    return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
  }

  const { URLs } = await req.json();
  if (!URLs || !Array.isArray(URLs)) {
    return NextResponse.json({ error: "Invalid input: URLs should be an array" }, { status: 400 });
  }
  // const URLs = formData.getAll("URLs") as string[];

  // Process each URL individually

  for (const url of URLs) {
    await fetchWithRefresh(`${API_URL}/image/?url=${url}`, {
      method: "DELETE",
    });
  }

  // Return combined results
  return NextResponse.json({}, { status: 200 });
}
