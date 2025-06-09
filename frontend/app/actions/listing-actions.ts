'use server';

import { deleteListing } from '@/app/lib/api/listing/requests';
import { revalidatePath } from 'next/cache';

export async function deleteListingAction(listingId: string) {
  try {
    const response = await deleteListing(listingId);
    console.log('delete', response);
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to delete listing' };
  }
}

export async function revalidateListingsAction() {
  revalidatePath('/account/listings');
}
