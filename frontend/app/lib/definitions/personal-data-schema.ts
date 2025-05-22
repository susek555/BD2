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
