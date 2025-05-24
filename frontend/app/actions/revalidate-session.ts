'use server';

import { revalidatePath } from 'next/cache';

export async function revalidateSesion(path: string, type: 'layout' | 'page') {
  revalidatePath(path, type);
}
