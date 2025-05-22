import { z } from 'zod';
import { PersonalDataSchema } from './personal-data-schema';

export const SignupFormSchema = PersonalDataSchema.and(
  z
    .object({
      password: z
        .string()
        .min(2, { message: 'Must be at least 2 characters long' })
        // .regex(/[a-zA-Z]/, { message: 'Contain at least one letter.' })
        // .regex(/[0-9]/, { message: 'Contain at least one number.' })
        // .regex(/[^a-zA-Z0-9]/, {
        //   message: 'Contain at least one special character.',
        // })
        .trim(),
      confirm_password: z.string().trim(),
    })
    .refine((data) => data.password === data.confirm_password, {
      message: "Passwords don't match",
      path: ['confirm_password'],
    })
    .transform((data) => {
      const { confirm_password: _, ...rest } = data;
      return rest;
    }),
);

export type SignupForm = z.input<typeof SignupFormSchema>;

export type SignupFormState = {
  errors?: {
    username?: string[];
    email?: string[];
    password?: string[];
    confirm_password?: string[];
    selector?: string[];

    person_name?: string[];
    person_surname?: string[];

    company_name?: string[];
    company_nip?: string[];
  };
  values?: {
    selector?: string;
    username?: string;
    email?: string;

    person_name?: string;
    person_surname?: string;

    company_name?: string;
    company_nip?: string;
  };
};

export type LoginError = {
  credentials?: string[];
  server?: string[];
};
