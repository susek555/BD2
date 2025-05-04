import { z } from 'zod';

export const SignupFormSchema = z.object({
  selector: z.enum(['P', 'C']),
  username: z
    .string()
    .min(2, { message: 'Username must be at least 2 characters long.' })
    .trim(),
  email: z.string().email({ message: 'Please enter a valid email.' }).trim(),
  password: z
    .string()
    .min(2, { message: 'Be at least 2 characters long' })
    // .regex(/[a-zA-Z]/, { message: 'Contain at least one letter.' })
    // .regex(/[0-9]/, { message: 'Contain at least one number.' })
    // .regex(/[^a-zA-Z0-9]/, {
    //   message: 'Contain at least one special character.',
    // })
    .trim(),
  confirm_password: z.string().trim(),
})
  .refine(
    (data) => data.password === data.confirm_password,
    {
      message: "Passwords don't match",
      path: ["confirm_password"],
    }
  )
  .and(
    z.discriminatedUnion('selector', [
      z.object({
        selector: z.literal('P'),
        person_name: z.string().min(1, { message: 'Name is required' }),
        person_surname: z.string().min(1, { message: 'Surname is required' }),
      }),
      z.object({
        selector: z.literal('C'),
        company_name: z.string().min(1, { message: 'Business name is required' }),
        company_nip: z
          .string()
          .min(10, { message: 'Be 10 characters long' })
          .max(10, { message: 'Be 10 characters long' })
          .regex(/^\d+$/, { message: 'Contain only digits' }),
      })
    ])
  )
  .transform(data => {
    const { confirm_password: confir_password, ...rest } = data;
    return rest;
  });

export type SignupForm = z.input<typeof SignupFormSchema>;

export type SignupFormState = {
  errors?: {
    username?: string[]
    email?: string[]
    password?: string[]
    confirm_password?: string[]
    selector?: string[]

    person_name?: string[]
    person_surname?: string[]

    company_name?: string[]
    company_nip?: string[]
  }
  values?: {
    selector?: string
    username?: string
    email?: string

    person_name?: string
    person_surname?: string

    company_name?: string
    company_nip?: string
  }
}
