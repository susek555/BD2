export type PersonalDataFormState = {
  errors?: {
    username?: string[];
    email?: string[];

    personName?: string[];
    personSurname?: string[];

    companyName?: string[];
    companyNip?: string[];
  };
  values?: {
    username?: string;
    email?: string;

    personName?: string;
    personSurname?: string;

    companyName?: string;
    companyNip?: string;
  };
};
