
if ($args.count -eq 1) {
    if (test-path $args[0]) {
        $text = get-content $args[0]

        foreach ($row in $text) {
            $numbers = $row.split("_")
            $max = [math]::max($numbers[0], $numbers[1])
            $max = [math]::max($numbers[2], $max)

            $str = "$row  =>  $max"
            Add-content result.txt $str
        }
    } else {
        write-host "Invalid file"
    }
} else {
    write-host "Not good"
}

