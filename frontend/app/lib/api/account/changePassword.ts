import { PasswordValidationResult } from '../auth';

export async function changePassword({
  id,
  currentPassword,
  newPassword,
}: {
  id: number;
  currentPassword: string;
  newPassword: string;
}): Promise<PasswordValidationResult> {
  const response = await fetch('/api/account/update/password', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ id, current_password: currentPassword, new_password: newPassword }),
  });

  return {};
}
