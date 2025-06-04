import { PasswordValidationResult } from '../auth';

export async function changePassword({
  currentPassword,
  newPassword,
}: {
  currentPassword: string;
  newPassword: string;
}): Promise<PasswordValidationResult> {
  try {
    const response = await fetch('/api/account/update/password', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        old_password: currentPassword,
        new_password: newPassword,
      }),
    });

    console.log(response);

    if (!response.ok) {
      const errorData: PasswordValidationResult = await response.json();
      return errorData;
    }

    return {};
  } catch (error) {
    console.error('Error updating password:', error);
    return { errors: { other: ['Something went wrong'] } };
  }
}
