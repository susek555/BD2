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

  const formData = await req.formData();
  const URLs = formData.getAll("URLs") as string[];

  // Process each URL individually
  const results = [];

  for (const url of URLs) {
    const response = await fetchWithRefresh(`${API_URL}/image/?url=${url}`, {
      method: "DELETE",
    });

    const result = await response.json();
    results.push(result);
  }

  // Return combined results
  return NextResponse.json({ results }, { status: 200 });
}
