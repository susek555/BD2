import { PasswordData } from '@/app/lib/definitions/user';
import { NextRequest, NextResponse } from 'next/server';

export async function PUT(request: NextRequest) {
  try {
    const userData: PasswordData = await request.json();

    return NextResponse.json({ status: 200 });
  } catch (error) {
    console.error('Error changing password:', error);
    return NextResponse.json(
      { errors: { other: ['Failed to change password'] } },
      { status: 500 },
    );
  }
}
