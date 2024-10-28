import { Box, Button, TextField } from '@mui/material'
import { Control, Controller } from 'react-hook-form'
import { RedirectURIsEdit, RedirectURIsFields } from './RedirectURIsEdit'

export type ClientInput = {
  name: string
  redirectUris: string[]
}
export type ClientFormProps = {
  submit: () => void
  control: Control<ClientInput>
  redirectURIsControl: Control<RedirectURIsFields>
}

export const ClientForm = ({
  submit,
  control,
  redirectURIsControl,
}: ClientFormProps): JSX.Element => {
  const requiredRule = { required: 'このフィールドは必須です。' }
  return (
    <Box onSubmit={submit} component="form" sx={{ '> *': { margin: 2 } }}>
      <Box>
        <Controller
          control={control}
          name="name"
          rules={requiredRule}
          render={({ field, fieldState }) => (
            <TextField
              fullWidth
              label="Name"
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              {...field}
            />
          )}
        />
      </Box>
      <RedirectURIsEdit control={redirectURIsControl} />
      <Box>
        <Button type="submit" fullWidth variant="contained">
          Submit
        </Button>
      </Box>
    </Box>
  )
}
