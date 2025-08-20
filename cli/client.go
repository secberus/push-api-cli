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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/secberus/push-api-cli/model"

	api "github.com/secberus/go-push-api/api/v1"
	service "github.com/secberus/go-push-api/service/v1/push"
	v1 "github.com/secberus/go-push-api/types/v1"
)

const (
	FormatUnknown = 0
	FormatJSON    = 1
	FormatYAML    = 2
	FormatCSV     = 3
)

var (
	UseJsonFormat = false
	UseYamlFormat = false
	UseCsvFormat  = false
)

type (
	csvSyncDataRow struct {
	}

	csvUpsertRecordsRow struct {
	}

	csvDeleteRecordsRow struct {
		TableName  string
		ColumnName string
		PrimaryKey bool
	}
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
		//		newGetIndexCommand(client),
		//		newListIndexesCommand(client),
		//		newCreateIndexCommand(client),
		//		newDropIndexCommand(client),
		//		newAlterTableCommand(client),
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
	format := FormatJSON

	cmd := &cobra.Command{
		Use:   "create-table [ --yaml | --json | --csv ] [ -f <filename> ]",
		Short: "create table",
		Long:  "create table",
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			input, err := readFileOrStdin(command)
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			switch {
			case UseCsvFormat:
				format = FormatCSV
			case UseYamlFormat:
				format = FormatYAML
			default:
				format = FormatJSON
			}

			var table model.Table
			if err := decodeInput(input, format, &table); err != nil {
				return fmt.Errorf("failed to decode input: %w", err)
			}

			_, err = client.CreateTable(ctx, &api.CreateTableInput{
				Table: &v1.Table{
					Name:     table.Name,
					SyncType: resolveSyncType(table.SyncType),
					Columns:  encodeColumns(table.Columns),
				},
			})
			if err != nil {
				return fmt.Errorf("failed to create table: %w", err)
			}

			fmt.Fprintf(command.OutOrStdout(), "created table %s\n", table.Name)
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("file", "f", "", "Path to the file containing the table schema.")
	cmd.Flags().BoolVar(&UseJsonFormat, "json", false, "Parse the input as JSON.  This is the default if no format is specified.")
	cmd.Flags().BoolVar(&UseYamlFormat, "yaml", false, "Parse the input as YAML.")
	cmd.Flags().BoolVar(&UseCsvFormat, "csv", false, "Parse the input as CSV.")
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
	format := FormatJSON

	cmd := &cobra.Command{
		Use:   "upsert-records [ --json | --yaml | --csv ] -f <filename> | <input>",
		Short: "upsert records",
		Long:  "upsert records",
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			input, err := readFileOrStdin(command)
			if err != nil {
				return err
			}

			switch {
			case UseCsvFormat:
				format = FormatCSV
			case UseYamlFormat:
				format = FormatYAML
			default:
				format = FormatJSON
			}

			var records []model.Record
			if err := decodeInput(input, format, &records); err != nil {
				return err
			}
			_, err = client.UpsertRecords(ctx, &api.UpsertRecordsInput{
				Records: encodeRecords(records),
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(command.OutOrStdout(), "upserted %d records", len(records))
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("file", "f", "", "Path to the file containing the records to upsert")
	cmd.Flags().BoolVar(&UseJsonFormat, "json", false, "Parse the input as JSON.  This is the default if no format is specified.")
	cmd.Flags().BoolVar(&UseYamlFormat, "yaml", false, "Parse the input as YAML.")
	cmd.Flags().BoolVar(&UseCsvFormat, "csv", false, "Parse the input as CSV.")
	cmd.MarkFlagFilename("file")
	return cmd
}

func newDeleteRecordsCommand(client service.PushServiceClient) *cobra.Command {
	format := FormatJSON

	cmd := &cobra.Command{
		Use:   "delete-records [ --json | --yaml | --csv ] -f <filename> | <input>",
		Short: "delete records",
		Long:  "delete records",
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			input, err := readFileOrStdin(command)
			if err != nil {
				return err
			}

			switch {
			case UseCsvFormat:
				format = FormatCSV
			case UseYamlFormat:
				format = FormatYAML
			default:
				format = FormatJSON
			}

			var records []model.Record
			if err := decodeInput(input, format, &records); err != nil {
				return err
			}
			_, err = client.DeleteRecords(ctx, &api.DeleteRecordsInput{
				PrimaryKey: encodeRecords(records),
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(command.OutOrStdout(), "deleted %d records", len(records))
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("file", "f", "", "Path to the file containing the records to delete")
	cmd.Flags().BoolVar(&UseJsonFormat, "json", false, "Parse the input as JSON.  This is the default if no format is specified.")
	cmd.Flags().BoolVar(&UseYamlFormat, "yaml", false, "Parse the input as YAML.")
	cmd.Flags().BoolVar(&UseCsvFormat, "csv", false, "Parse the input as CSV.")
	cmd.MarkFlagFilename("file")
	return cmd
}

func newSyncDataCommand(client service.PushServiceClient) *cobra.Command {
	format := FormatJSON

	cmd := &cobra.Command{
		Use:   "sync-data [ --json | --yaml | --csv ] -f <filename> | <input>",
		Short: "sync data",
		Long:  "sync data",
		RunE: func(command *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			input, err := readFileOrStdin(command)
			if err != nil {
				return err
			}

			switch {
			case UseCsvFormat:
				format = FormatCSV
			case UseYamlFormat:
				format = FormatYAML
			default:
				format = FormatJSON
			}

			var data []*api.SyncDataInput
			if err := decodeInput(input, format, &data); err != nil {
				return err
			}

			stream, err := client.SyncData(ctx)
			if err != nil {
				return err
			}

			recordCount := 0
			for _, d := range data {
				if err := stream.Send(d); err != nil {
					return err
				}
				recordCount += len(d.Records)
			}

			if err := stream.CloseSend(); err != nil {
				return err
			}
			fmt.Fprintf(command.OutOrStdout(), "synced %d records", recordCount)
			return nil
		},
		SilenceUsage: true,
	}
	cmd.Flags().StringP("file", "f", "", "Path to the file containing the data to sync")
	cmd.Flags().BoolVar(&UseJsonFormat, "json", false, "Parse the input as JSON.  This is the default if no format is specified.")
	cmd.Flags().BoolVar(&UseYamlFormat, "yaml", false, "Parse the input as YAML.")
	cmd.Flags().BoolVar(&UseCsvFormat, "csv", false, "Parse the input as CSV.")
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

// detects the format of the input from the io.,Reader and attempts to decode into the v paramter.
func decodeInput(r io.Reader, format int, v any) error {

	switch format {
	case FormatJSON:
		if err := decodeJSON(r, v); err != nil {
			return fmt.Errorf("failed to decode as JSON: %w", err)
		}
	case FormatYAML:
		if err := decodeYAML(r, v); err != nil {
			return fmt.Errorf("failed to decode as YAML: %w", err)
		}
	case FormatCSV:
		if err := decodeCSV(r, v); err != nil {
			return fmt.Errorf("failed to decode as CSV: %w", err)
		}
	default:
		return fmt.Errorf("cannot decode unknown input format")
	}
	return nil
}

func decodeJSON(input io.Reader, v any) error {
	decoder := json.NewDecoder(input)
	return decoder.Decode(v)
}

func decodeYAML(input io.Reader, v any) error {
	decoder := yaml.NewDecoder(input)
	return decoder.Decode(v)
}

func decodeCSV(input io.Reader, v any) error {
	return nil
}

func resolveSyncType(s string) v1.TableSyncType {

	if strings.Contains(strings.ToLower(s), "truncate") {
		return v1.TableSyncType_TABLE_SYNC_TYPE_TRUNCATE
	} else if strings.Contains(strings.ToLower(s), "append") {
		return v1.TableSyncType_TABLE_SYNC_TYPE_APPEND
	}

	return v1.TableSyncType_TABLE_SYNC_TYPE_UNSPECIFIED
}

func encodeRecords(r []model.Record) []*v1.Record {

	records := make([]*v1.Record, len(r))
	for x, rec := range r {
		records[x] = &v1.Record{
			TableName: rec.TableName,
			Columns:   encodeColumns(rec.Columns),
		}
	}
	return records
}

func encodeColumns(c []model.Column) []*v1.Column {

	cols := make([]*v1.Column, len(c))
	for x, col := range c {

		dt := &v1.DataType{}

		switch {
		case col.DataType.Boolean != nil:
			dt.Union = &v1.DataType_Boolean{
				Boolean: &v1.Boolean{
					Value:   col.DataType.Boolean.Value,
					Default: col.DataType.Boolean.Default,
				},
			}
		case col.DataType.Integer != nil:
			dt.Union = &v1.DataType_Integer{
				Integer: &v1.Integer{
					Value:   col.DataType.Integer.Value,
					Default: col.DataType.Integer.Default,
				},
			}
		case col.DataType.Text != nil:
			dt.Union = &v1.DataType_Text{
				Text: &v1.Text{
					Value:   col.DataType.Text.Value,
					Default: col.DataType.Text.Default,
				},
			}
		}

		cols[x] = &v1.Column{
			Name:       col.Name,
			Nillable:   col.Nillable,
			Unique:     col.Unique,
			PrimaryKey: col.PrimaryKey,
			DataType:   dt,
		}
	}
	return cols
}
