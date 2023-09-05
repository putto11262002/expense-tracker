import axios from "axios";


export function extractErrorMessage(err: unknown) {
  let errMessage: string;

  if (axios.isAxiosError<{ error: string }>(err) && err.response) {
    errMessage = err.response.data.error;
  } else {
    console.error(err);
    errMessage = "Something went wrong please try again";
  }

  return Promise.resolve(errMessage);
}
