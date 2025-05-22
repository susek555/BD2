import { UserProfile } from '@/app/lib/definitions/user';
import { NextRequest, NextResponse } from 'next/server';

export async function PUT(request: NextRequest) {
  try {
    const userData: UserProfile = await request.json();

    return NextResponse.json({ status: 200 });
  } catch (error) {
    console.error('Error updating profile:', error);
    return NextResponse.json(
      { errors: { other: ['Failed to update profile'] } },
      { status: 500 },
    );
  }
}
