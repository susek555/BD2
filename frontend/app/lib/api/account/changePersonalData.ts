import { UserProfile } from '../../definitions/user';
import { AccountFieldValidationResult } from '../auth';

export async function changePersonalData(
  newUserData: UserProfile,
): Promise<AccountFieldValidationResult> {
  try {
    const response = await fetch('/api/account/update/profile', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(newUserData),
    });

    if (!response.ok) {
      const errorData: AccountFieldValidationResult = await response.json();
      return errorData;
    }

    return {};
  } catch (error) {
    console.error('Error updating profile:', error);
    return { errors: { other: ['Network error occurred'] } };
  }
}
