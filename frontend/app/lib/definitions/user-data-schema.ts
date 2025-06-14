import { z } from 'zod';

export const PersonalDataSchema = z
  .object({
    selector: z.enum(['P', 'C']),
    username: z
      .string()
      .min(2, { message: 'Username must be at least 2 characters long.' })
      .trim(),
    email: z.string().email({ message: 'Please enter a valid email.' }).trim(),
  })
  .and(
    z.discriminatedUnion('selector', [
      z.object({
        selector: z.literal('P'),
        person_name: z.string().min(1, { message: 'Name is required' }),
        person_surname: z.string().min(1, { message: 'Surname is required' }),
      }),
      z.object({
        selector: z.literal('C'),
        company_name: z
          .string()
          .min(1, { message: 'Company name is required' }),
        company_nip: z
          .string()
          .min(10, { message: 'Must be 10 characters long' })
          .max(10, { message: 'Must be 10 characters long' })
          .regex(/^\d+$/, { message: 'Can contain only digits' }),
      }),
    ]),
  );

export const ChangePasswordSchema = z
  .object({
    current_password: z.string(),
    new_password: z
      .string()
      .min(2, { message: 'Must be at least 2 characters long' })
      // .regex(/[a-zA-Z]/, { message: 'Contain at least one letter.' })
      // .regex(/[0-9]/, { message: 'Contain at least one number.' })
      // .regex(/[^a-zA-Z0-9]/, {
      //   message: 'Contain at least one special character.',
      // })
      .trim(),
    confirm_new_password: z.string().trim(),
  })
  .refine((data) => data.new_password === data.confirm_new_password, {
    message: "Passwords don't match",
    path: ['confirm_new_password'],
  })
  .transform((data) => {
    const { confirm_new_password: _, ...rest } = data;
    return rest;
  });
