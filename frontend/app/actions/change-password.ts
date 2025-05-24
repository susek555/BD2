import { changePassword } from '../lib/api/account/changePassword';
import { ChangePasswordFormState } from '../lib/definitions/account-settings';
import { ChangePasswordSchema } from '../lib/definitions/user-data-schema';

export async function changePasswordAction(
  state: ChangePasswordFormState,
  formData: FormData,
  id: number,
): Promise<ChangePasswordFormState> {
  const formDataObj = Object.fromEntries(formData.entries());
  const validatedFields = ChangePasswordSchema.safeParse(formDataObj);

  if (!validatedFields.success) {
    console.log(
      'Password change validation errors:',
      validatedFields.error.flatten().fieldErrors,
    );
    return {
      errors: validatedFields.error.flatten().fieldErrors,
    };
  }

  const response = await changePassword({
    id,
    currentPassword: validatedFields.data.current_password,
    newPassword: validatedFields.data.new_password,
  });

  if (response.errors) {
    return { errors: response.errors };
  } else {
    return { success: true };
  }
}
