import camelcaseKeys from 'camelcase-keys';
import { updatePersonalData } from '../lib/api/account/updatePersonalData';
import { AccountFieldValidationResult } from '../lib/api/auth';
import { PersonalDataFormState } from '../lib/definitions/account-settings';
import { PersonalDataSchema } from '../lib/definitions/personal-data-schema';
import { UserProfile } from '../lib/definitions/user';

export async function updateProfile(
  state: PersonalDataFormState,
  formData: FormData,
  id: number,
): Promise<{ state: PersonalDataFormState; newUserData?: UserProfile }> {
  const formDataObj = Object.fromEntries(formData.entries());
  const validatedFields = PersonalDataSchema.safeParse(formDataObj);

  if (!validatedFields.success) {
    console.log(
      'Profile update validation errors:',
      validatedFields.error.flatten().fieldErrors,
    );
    return {
      state: {
        errors: validatedFields.error.flatten().fieldErrors,
        values: formDataObj.values as PersonalDataFormState['values'],
      },
    };
  }
  const updateResult: AccountFieldValidationResult = await updatePersonalData({
    id,
    ...validatedFields.data,
  });

  if (updateResult.errors) {
    return {
      state: {
        errors: updateResult.errors,
        values: Object.fromEntries(
          Object.entries(formDataObj).filter(
            ([key]) => !key.includes('password'),
          ),
        ) as PersonalDataFormState['values'],
      },
    };
  } else {
    const updatedUserProfile: UserProfile = {
      id,
      ...camelcaseKeys(validatedFields.data, { deep: true }),
    };
    return { state: {}, newUserData: updatedUserProfile };
  }
}
