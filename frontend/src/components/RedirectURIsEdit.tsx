'use client'
import {
  Box,
  Button,
  IconButton,
  List,
  ListItem,
  TextField,
  Typography,
} from '@mui/material'
import {
  Control,
  Controller,
  FieldValues,
  useFieldArray,
} from 'react-hook-form'
import DeleteIcon from '@mui/icons-material/Delete'
import { JSX } from 'react'

export interface RedirectURIsFields extends FieldValues {
  redirectURIs: Array<{ url: string }>
}
export type RedirectURIsEditProps = {
  control: Control<RedirectURIsFields>
}
export const RedirectURIsEdit = ({
  control,
}: RedirectURIsEditProps): JSX.Element => {
  const { fields, append, remove } = useFieldArray({
    control,
    name: 'redirectURIs',
  })

  return (
    <Box>
      <Typography variant="h5" component="h2">
        RedirectURIs
      </Typography>
      <List>
        {fields.map((item, index) => (
          <ListItem sx={{ display: 'flex' }} key={item.id}>
            <Controller
              control={control}
              name={`redirectURIs.${index}.url`}
              render={({ field, fieldState }) => (
                <TextField
                  fullWidth
                  label="Redirect URI"
                  error={fieldState.invalid}
                  helperText={fieldState.error?.message}
                  {...field}
                />
              )}
            />
            <Box sx={{ ml: 1 }}>
              <IconButton onClick={() => remove(index)}>
                <DeleteIcon fontSize="large" />
              </IconButton>
            </Box>
          </ListItem>
        ))}
      </List>
      <Box sx={{ pl: 2 }}>
        <Button onClick={() => append({ url: '' })}>Add Redirect URI</Button>
      </Box>
    </Box>
  )
}
