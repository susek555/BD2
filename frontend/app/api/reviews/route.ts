import { NewReview, UpdatedReview } from '@/app/lib/definitions/reviews';
import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const body: NewReview = await request.json();

    return NextResponse.json(
      { message: 'Created successfully' },
      { status: 201 },
    );
  } catch (error) {
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 },
    );
  }
}

export async function PUT(request: NextRequest) {
  try {
    const body: UpdatedReview = await request.json();

    return NextResponse.json(
      { message: 'Updated successfully' },
      { status: 200 },
    );
  } catch (error) {
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 },
    );
  }
}
