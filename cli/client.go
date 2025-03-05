/*
 * Copyright 2024 Secberus, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package cli

import (
  "bytes"
	"context"
  "encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
  "gopkg.in/yaml.v3"

	api "github.com/secberus/go-push-api/api/v1"
  types "github.com/secberus/go-push-api/types/v1"
	service "github.com/secberus/go-push-api/service/v1/push"
)

func New(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use: "push-api-cli <command>",
	}
	cmd.AddCommand(
		newListTablesCommand(client),
		newGetTableCommand(client),
		newCreateTableCommand(client),
		newDropTableCommand(client),
		newTruncateTableCommand(client),
		newGetIndexCommand(client),
		newListIndexesCommand(client),
		newCreateIndexCommand(client),
		newDropIndexCommand(client),
		newAlterTableCommand(client),
		newUpsertRecordsCommand(client),
		newDeleteRecordsCommand(client),
		newSyncDataCommand(client),
	)
	cmd.Version = "1.0"
	return cmd
}

func newListTablesCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-tables",
		Short: "list tables",
		Long:  "list tables",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      stream, err := client.ListTables(ctx, &api.ListTablesInput{})
      if err != nil {
        return err
      }

      for {
        output, err := stream.Recv()
        if err == io.EOF {
          break
        }
        if err != nil {
          return err
        }

        table := output.GetTable()
        if table != nil {
          fmt.Fprintf(command.OutOrStdout(), "%s\n", table.GetName())
        }
      }
			return nil
		},
		SilenceUsage: true,
	}
	return cmd
}

func newGetTableCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-table <name>",
		Short: "get table",
		Long:  "get table",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      response, err := client.GetTable(ctx, &api.GetTableInput{
        TableName: args[0],
      })
      if err != nil {
        return err
      }

      table := response.GetTable()
      if table != nil {
        y, err := yaml.Marshal(table)
        if err != nil {
          return err
        }

        fmt.Fprintf(command.OutOrStdout(), "%s\n", string(y))
      }

			return nil
		},
		SilenceUsage: true,
	}
	return cmd
}

func newCreateTableCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-table [ -f <filename> ]",
		Short: "create table",
		Long:  "create table",
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      input, err := readFileOrStdin(command)
      if err != nil {
        return err
      }

      var table types.Table
      if err := decodeYAMLOrJSON(input, &table); err != nil {
        return err
      }
      _, err = client.CreateTable(ctx, &api.CreateTableInput{
        Table: &table,
      })
      if err != nil {
        return err
      }

      fmt.Fprintf(command.OutOrStdout(), "created table %s\n", table.Name)
			return nil
		},
		SilenceUsage: true,
	}
  cmd.Flags().StringP("file", "f", "", "Path to the file containing the table schema")
  cmd.MarkFlagFilename("file")
  return cmd
}

func newDropTableCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drop-table <name>",
		Short: "drop table",
		Long:  "drop table",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      _, err := client.DropTable(ctx, &api.DropTableInput{
        TableName: args[0],
      })
      if err != nil {
        return err
      }

      fmt.Fprintf(command.OutOrStdout(), "dropped table %s\n", args[0])
			return nil
		},
		SilenceUsage: true,
	}
  return cmd
}

func newAlterTableCommand(client service.PushServiceClient) *cobra.Command {
  return nil
}

func newTruncateTableCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "truncate-table <name>",
		Short: "truncate table",
		Long:  "truncate table",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      _, err := client.TruncateTable(ctx, &api.TruncateTableInput{
        TableName: args[0],
      })
      if err != nil {
        return err
      }

      fmt.Fprintf(command.OutOrStdout(), "truncated table %s\n", args[0])
			return nil
		},
		SilenceUsage: true,
	}
  return cmd
}

func newCreateIndexCommand(client service.PushServiceClient) *cobra.Command {
  return nil
}

func newGetIndexCommand(client service.PushServiceClient) *cobra.Command {
  return nil
}

func newDropIndexCommand(client service.PushServiceClient) *cobra.Command {
  return nil
}

func newListIndexesCommand(client service.PushServiceClient) *cobra.Command {
  return nil
}

func newUpsertRecordsCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-records -f <filename> | <input>",
		Short: "upsert records",
		Long:  "upsert records",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      response, err := client.GetTable(ctx, &api.GetTableInput{
        TableName: args[0],
      })
      if err != nil {
        return err
      }

      table := response.GetTable()
      if table != nil {
        fmt.Printf("%s\n", table.GetName())
      }

			return nil
		},
		SilenceUsage: true,
	}
  cmd.Flags().StringP("file", "f", "", "Path to the file containing the records to upsert")
  cmd.MarkFlagFilename("file")
  return cmd
}

func newDeleteRecordsCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-records -f <filename> | <input>",
		Short: "delete records",
		Long:  "delete records",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      response, err := client.GetTable(ctx, &api.GetTableInput{
        TableName: args[0],
      })
      if err != nil {
        return err
      }

      table := response.GetTable()
      if table != nil {
        fmt.Printf("%s\n", table.GetName())
      }

			return nil
		},
		SilenceUsage: true,
	}
  cmd.Flags().StringP("file", "f", "", "Path to the file containing the records to delete")
  cmd.MarkFlagFilename("file")
  return cmd
}

func newSyncDataCommand(client service.PushServiceClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-data -f <filename> | <input>",
		Short: "sync data",
		Long:  "sync data",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

      response, err := client.GetTable(ctx, &api.GetTableInput{
        TableName: args[0],
      })
      if err != nil {
        return err
      }

      table := response.GetTable()
      if table != nil {
        fmt.Printf("%s\n", table.GetName())
      }

			return nil
		},
		SilenceUsage: true,
	}
  cmd.Flags().StringP("file", "f", "", "Path to the file containing the data to sync")
  cmd.MarkFlagFilename("file")
  return cmd
}

func readFileOrStdin(command *cobra.Command) (io.Reader, error) {
  // set reader to stdin if no file is given
  var inputReader io.Reader = command.InOrStdin()

  filename, err := command.Flags().GetString("file")
  if err != nil {
    return nil, fmt.Errorf("failed to parse argument: %v", err)
  }

  if filename != "" && filename != "-" {
    file, err := os.Open(filename)
    if err != nil {
      return nil, fmt.Errorf("failed to open file: %v", err)
    }
    inputReader = file
  }

  return inputReader, nil
}

func decodeYAMLOrJSON(input io.Reader, v any) error {
	// reads the entire content from the reader into a buffer.  this'll need to change eventually to be able to handle larger inputs.
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, input)
	if err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	data := buf.Bytes()
	// check if the content starts with JSON
	if json.Valid(data) {
		// If valid JSON, decode it
		decoder := json.NewDecoder(bytes.NewReader(data))
		return decoder.Decode(v)
	}

	// If it's not valid JSON, try to decode as YAML
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(v)
}

