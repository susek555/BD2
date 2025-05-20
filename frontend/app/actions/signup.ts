import { registerResult, registerUser } from "@/app/lib/api/auth";
import { SignupFormSchema, SignupFormState } from "@/app/lib/definitions/reviews";
import { permanentRedirect } from "next/navigation";


export async function signup(
  state: SignupFormState,
  formData: FormData
): Promise<SignupFormState> {
  console.log("Signup form data:", Object.fromEntries(formData.entries()));

  const formDataObj = Object.fromEntries(formData.entries());
  const validatedFields = SignupFormSchema.safeParse(formDataObj);
  console.log("Signup validation result:", validatedFields);

  if (!validatedFields.success) {
    console.log("Signup validation errors:", validatedFields.error.flatten().fieldErrors);
    return {
      errors: validatedFields.error.flatten().fieldErrors,
      values: Object.fromEntries(
        Object.entries(formDataObj).filter(([key]) => !key.includes('password'))
      ) as SignupFormState['values']
    };
  }

  const signupResult: registerResult = await registerUser(validatedFields.data);

  if (signupResult.errors) {
    return {
      errors: signupResult.errors,
      values: Object.fromEntries(
        Object.entries(formDataObj).filter(([key]) => !key.includes('password'))
      ) as SignupFormState['values']
    }
  } else {
    permanentRedirect("/")
  }
}
