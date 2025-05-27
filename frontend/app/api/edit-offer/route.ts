// app/api/edit-offer/route.ts
import { fetchWithRefresh } from "@/app/lib/api/fetchWithRefresh";
import { getServerSession } from "next-auth";
import { authConfig } from "@/app/lib/authConfig";
import { NextRequest, NextResponse } from "next/server";

const API_URL = process.env.API_URL;

export async function PUT(
    req: NextRequest
) {
    const session = await getServerSession(authConfig);
    if (!session) {
        return NextResponse.json({ error: "Not authenticated" }, { status: 401 });
    }

    const data = await req.json();
    const id = req.nextUrl.pathname.split('/').pop();

    const response = await fetchWithRefresh(`${API_URL}/sale-offer/${id}`, {
        // TODO update URL
        method: "PUT",
        body: JSON.stringify(data),
        headers: {
            "Content-Type": "application/json",
        },
    });

    const result = await response.json();
    return NextResponse.json(result, { status: response.status });
}