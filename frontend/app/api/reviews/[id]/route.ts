import { NextResponse } from 'next/server';

export async function DELETE(
  { params }: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await params;

    return NextResponse.json(
      { message: 'Deleted successfully' },
      { status: 200 },
    );
  } catch (error) {
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 },
    );
  }
}
