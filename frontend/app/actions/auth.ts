import { SignupFormSchema, SignupFormState } from "@/app/lib/definitions";

export async function signup(
  state: SignupFormState,
  formData: FormData
): Promise<SignupFormState> {
  console.log("Form data:", Object.fromEntries(formData.entries()));

  const formDataObj = Object.fromEntries(formData.entries());

  const validatedFields = SignupFormSchema.safeParse(formDataObj);
  console.log("Validation result:", validatedFields);

  if (!validatedFields.success) {
    console.log("Validation errors:", validatedFields.error.flatten().fieldErrors);
    return {
      errors: validatedFields.error.flatten().fieldErrors,
      values: formDataObj as SignupFormState['values']
    };
  }

  return {}
}
