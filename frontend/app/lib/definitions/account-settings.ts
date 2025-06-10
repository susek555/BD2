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

export type ChangePasswordFormState = {
  errors?: {
    current_password?: string[];
    new_password?: string[];
    confirm_new_password?: string[];
    other?: string[];
  };
  success?: boolean;
};
